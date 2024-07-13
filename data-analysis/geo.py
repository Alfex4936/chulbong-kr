import json
import os
import time

import folium
import geopandas as gpd
import pandas as pd
from esda.getisord import G, G_Local
from esda.moran import Moran
from folium.plugins import HeatMap
from geopy.exc import GeocoderServiceError, GeocoderTimedOut
from geopy.geocoders import Nominatim
from libpysal.weights import Queen
from scipy.spatial import cKDTree
from shapely.geometry import Point
from sklearn.cluster import DBSCAN


# Function to standardize province names
def standardize_province(province):
    mapping = {
        "경기": "경기도",
        "경기도": "경기도",
        "서울": "서울특별시",
        "서울특별시": "서울특별시",
        "부산": "부산광역시",
        "부산광역시": "부산광역시",
        "대구": "대구광역시",
        "대구광역시": "대구광역시",
        "인천": "인천광역시",
        "인천광역시": "인천광역시",
        "제주": "제주특별자치도",
        "제주특별자치도": "제주특별자치도",
        "제주도": "제주특별자치도",
        "대전": "대전광역시",
        "대전광역시": "대전광역시",
        "울산": "울산광역시",
        "울산광역시": "울산광역시",
        "광주": "광주광역시",
        "광주광역시": "광주광역시",
        "세종": "세종특별자치시",
        "세종특별자치시": "세종특별자치시",
        "강원": "강원특별자치도",
        "강원도": "강원특별자치도",
        "강원특별자치도": "강원특별자치도",
        "경남": "경상남도",
        "경상남도": "경상남도",
        "경북": "경상북도",
        "경상북도": "경상북도",
        "전북": "전북특별자치도",
        "전라북도": "전북특별자치도",
        "충남": "충청남도",
        "충청남도": "충청남도",
        "충북": "충청북도",
        "충청북도": "충청북도",
        "전남": "전라남도",
        "전라남도": "전라남도",
    }
    return mapping.get(province, province)


geojson_file = "pullup_bars.geojson"

# Check if the geojson file already exists
if os.path.exists(geojson_file):
    # Load the GeoDataFrame from the existing file
    geo_df = gpd.read_file(geojson_file)
    print("Loaded geocoded data from existing GeoJSON file.")
else:
    # Load the data
    with open("markers.json", "r", encoding="utf-8") as file:
        data = json.load(file)

    # Extract addresses
    addresses = [entry["address"] for entry in data]

    # Initialize the geolocator with increased timeout
    geolocator = Nominatim(user_agent="chulbong_kr", timeout=10)
    location_data = []

    # Function to handle geocoding with retry on timeout
    def geocode_address(address, retries=3):
        for _ in range(retries):
            try:
                return geolocator.geocode(address)
            except (GeocoderTimedOut, GeocoderServiceError):
                time.sleep(1)  # Wait for a second before retrying
        print(f"Failed to geocode address: {address}")
        return None

    # Batch process addresses to handle rate limiting
    batch_size = 10

    # Create a progress bar
    with tqdm(total=len(addresses), desc="Geocoding addresses") as pbar:
        for i in range(0, len(addresses), batch_size):
            batch = addresses[i : i + batch_size]
            for address in batch:
                location = geocode_address(address)
                if location:
                    location_data.append(
                        (address, location.latitude, location.longitude)
                    )
                pbar.update(1)
            # Introduce a delay to handle rate limits
            time.sleep(1)

    # Create a DataFrame
    df = pd.DataFrame(location_data, columns=["address", "latitude", "longitude"])

    # Create a GeoDataFrame
    geometry = [Point(xy) for xy in zip(df["longitude"], df["latitude"])]
    geo_df = gpd.GeoDataFrame(df, geometry=geometry)

    # Assign a default CRS to avoid warnings
    geo_df.set_crs(epsg=4326, inplace=True)

    # Save the GeoDataFrame to a file
    geo_df.to_file(geojson_file, driver="GeoJSON")
    print(f"Geocoded data saved to {geojson_file}")

# Standardize province names
geo_df["province"] = geo_df["address"].apply(
    lambda x: standardize_province(x.split()[0])
)

# Visualize Distribution
# Create a map centered on South Korea
m = folium.Map(location=[36.5, 127.5], zoom_start=7)

# Add pull-up bar locations to the map
for idx, row in geo_df.iterrows():
    folium.Marker([row["latitude"], row["longitude"]], popup=row["address"]).add_to(m)

# Save the map
m.save("pullup_bars_map.html")
print("Distribution map saved to pullup_bars_map.html")

# Regional Distribution Analysis
region_counts = geo_df["province"].value_counts()
print(region_counts)

# Create a regional distribution map
m = folium.Map(location=[36.5, 127.5], zoom_start=7)
for province, count in region_counts.items():
    folium.CircleMarker(
        location=[
            geo_df.loc[geo_df["province"] == province, "latitude"].mean(),
            geo_df.loc[geo_df["province"] == province, "longitude"].mean(),
        ],
        radius=count / 10,
        popup=f"{province}: {count} pull-up bars",
        color="blue",
        fill=True,
        fill_color="blue",
    ).add_to(m)
m.save("regional_distribution_map.html")
print("Regional distribution map saved to regional_distribution_map.html")

# Clustering Analysis
coords = geo_df[["latitude", "longitude"]].values
db = DBSCAN(eps=0.01, min_samples=5).fit(coords)
geo_df["cluster"] = db.labels_

# Plot clusters
m = folium.Map(location=[36.5, 127.5], zoom_start=7)
for idx, row in geo_df.iterrows():
    color = "red" if row["cluster"] == -1 else "green"
    folium.CircleMarker(
        location=[row["latitude"], row["longitude"]],
        radius=5,
        popup=f"Cluster {row['cluster']}",
        color=color,
        fill=True,
        fill_color=color,
    ).add_to(m)
m.save("clustering_map.html")
print("Clustering map saved to clustering_map.html")

# Spatial Autocorrelation Analysis
w = Queen.from_dataframe(geo_df)
moran = Moran(geo_df["geometry"].apply(lambda x: x.x + x.y), w)
print(f"Moran's I: {moran.I}, p-value: {moran.p_sim}")

# Hot Spot Analysis (Getis-Ord Gi*)
# Calculate Local G statistics (Getis-Ord Gi*)
g_local = G_Local(geo_df["geometry"].apply(lambda x: x.x + x.y), w)
geo_df["GiZ"] = g_local.Zs

# Plot hot spots
m = folium.Map(location=[36.5, 127.5], zoom_start=7)
for idx, row in geo_df.iterrows():
    if row["GiZ"] > 1.96:
        color = "red"  # Hot spot
        significance = "Hot Spot (High clustering)"
    elif row["GiZ"] < -1.96:
        color = "blue"  # Cold spot
        significance = "Cold Spot (Low clustering)"
    else:
        color = "green"  # Not significant
        significance = "Not Significant (Random distribution)"

    folium.CircleMarker(
        location=[row["latitude"], row["longitude"]],
        radius=5,
        popup=f"Z-Score: {row['GiZ']:.2f}<br>{significance}",
        color=color,
        fill=True,
        fill_color=color,
    ).add_to(m)

m.save("hotspot_map.html")
print("Hotspot map saved to hotspot_map.html")

# Density Heatmap
m = folium.Map(location=[36.5, 127.5], zoom_start=7)
heat_data = [[row["latitude"], row["longitude"]] for index, row in geo_df.iterrows()]
HeatMap(heat_data).add_to(m)
m.save("pullup_bars_heatmap.html")
print("Density heatmap saved to pullup_bars_heatmap.html")
