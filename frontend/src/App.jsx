import React, { useEffect, useState } from 'react';
import Fab from '@mui/material/Fab';
import MyLocationIcon from '@mui/icons-material/MyLocation';
import NavigationIcon from '@mui/icons-material/Navigation';

// https://apis.map.kakao.com/web/sample/multipleMarkerImage/

import './App.css';

function App() {
    const [map, setMap] = useState(null); // Store the map object in the state

    useEffect(() => {
        const script = document.createElement('script');
        script.src = "//dapi.kakao.com/v2/maps/sdk.js?appkey=dfdebaf84d7dda475fb8448c7d43c528&autoload=false";
        document.head.appendChild(script);

        script.onload = () => {
            window.kakao.maps.load(() => {
                var container = document.getElementById('map');
                var options = {
                    center: new window.kakao.maps.LatLng(37.52082051738358, 126.98824672418273),
                    level: 8
                };
                var mapCreated = new window.kakao.maps.Map(container, options);

                // Attach an event listener to the map object
                window.kakao.maps.event.addListener(mapCreated, 'click', function (mouseEvent) {
                    // Create a marker at the clicked position
                    var clickedPosition = mouseEvent.latLng;
                    var marker = new window.kakao.maps.Marker({
                        position: clickedPosition
                    });
                    marker.setMap(mapCreated);

                    // Optional: Display the latitude and longitude
                    console.log(`Latitude: ${clickedPosition.getLat()}, Longitude: ${clickedPosition.getLng()}`);
                });

                setMap(mapCreated);
            });
        };
    }, []);

    // Function to center map on user's current position
    const centerMapOnCurrentPosition = () => {
        if (map && navigator.geolocation) { // Check if map is loaded
            navigator.geolocation.getCurrentPosition((position) => {
                var newPos = new window.kakao.maps.LatLng(position.coords.latitude, position.coords.longitude);
                map.setLevel(3);
                map.setCenter(newPos);
            }, (error) => {
                console.error(error);
            });
        } else {
            alert('Geolocation is not supported by this browser or map is not loaded yet.');
        }
    };

    return (
        <div id="App">
            <div id="map" style={{ width: '100vw', height: '100vh' }}>
                {/* Render the Fab only if the map is loaded */}
                {map && (
                    <Fab
                        color="secondary" // Changed color for better visibility
                        aria-label="locate"
                        onClick={centerMapOnCurrentPosition}
                        sx={{
                            position: 'absolute',
                            bottom: 32,
                            right: 32,
                            color: 'white', // Text/icon color set to white for contrast
                            bgcolor: 'black', // Button color set to black for visibility
                            '&:hover': {
                                bgcolor: 'gray', // Button color on hover for better visibility
                            },
                            boxShadow: '0px 0px 10px rgba(0, 0, 0, 0.5)', // Shadow for 3D effect
                            border: '2px solid white', // Border to distinguish from the map
                        }}
                    >
                        <MyLocationIcon />
                    </Fab>

                )}
            </div>
        </div>
    );
}

export default App;
