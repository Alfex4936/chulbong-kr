"use client";

import GrowBox from "@/components/atom/GrowBox";
import { LocationIcon } from "@/components/icons/LocationIcons";
import SearchIcon from "@/components/icons/SearchIcon";
import { Input } from "@/components/ui/input";
import useInput from "@/hooks/common/useInput";
import useSearchLocationData from "@/hooks/query/useSearchLocationData";
import { useEffect, useState } from "react";
import useMapStore from "@/store/useMapStore";

const SearchInput = () => {
  const [query, setQuery] = useState("");
  const searchInput = useInput("");

  const { map } = useMapStore();

  const { data, isError } = useSearchLocationData(query);

  useEffect(() => {
    const handler = setTimeout(() => setQuery(searchInput.value), 300);

    return () => clearTimeout(handler);
  }, [searchInput.value]);

  const moveLocation = (x: number, y: number) => {
    const moveLatLon = new window.kakao.maps.LatLng(y, x);

    map?.setCenter(moveLatLon);
  };

  return (
    <div className="relative w-5/6 mx-auto mb-4">
      <div className="relative flex items-center justify-center">
        <div className="absolute top-1/2 left-2 -translate-y-1/2">
          <SearchIcon size={18} color="grey" />
        </div>
        <Input
          placeholder="장소, 위치를 입력하세요"
          value={searchInput.value}
          onChange={searchInput.handleChange}
          className="rounded-sm border border-solid border-grey placeholder:text-grey-dark pl-8"
        />
      </div>
      {searchInput.value.length > 0 && (
        <div className="absolute top-10 left-0 w-full z-10 bg-black rounded-sm border border-solid border-grey p-4">
          {isError && <div>잠시 후 다시 시도해 주세요.</div>}
          {data?.documents.length === 0 && <div>검색 결과가 없습니다.</div>}
          {data?.documents.map((document) => {
            return (
              <button
                key={document.id}
                className="flex w-full items-center p-1 px-2 rounded-sm mb-2 hover:bg-zinc-700"
                onClick={() =>
                  moveLocation(Number(document.x), Number(document.y))
                }
              >
                <div className="w-3/4">
                  <div className="truncate text-left mr-2">
                    {document.place_name}
                  </div>
                  <div className="truncate text-left text-xs text-grey-dark">
                    {document.address_name}
                  </div>
                </div>
                <GrowBox />
                <div>
                  <LocationIcon selected={false} size={18} />
                </div>
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default SearchInput;
