"use client";

import SearchIcon from "@/components/icons/SearchIcon";
import { Input } from "@/components/ui/input";
import useInput from "@/hooks/common/useInput";
import useSearchLocationData from "@/hooks/query/useSearchLocationData";
import { useEffect, useState } from "react";
import { ImCancelCircle } from "react-icons/im";
import MapSearchResult from "./MapSearchResult";

const MapSearch = () => {
  const [query, setQuery] = useState("");
  const searchInput = useInput("");

  const { data, isError } = useSearchLocationData(query);

  const [resultModal, setResultModal] = useState(false);

  useEffect(() => {
    const handler = setTimeout(() => setQuery(searchInput.value), 300);

    return () => clearTimeout(handler);
  }, [searchInput.value]);

  return (
    <div
      className={`absolute top-2 left-1/2 -translate-x-1/2 w-[90%] max-w-96 min-w-[280px] bg-black-light-2 z-50 rounded-sm`}
    >
      <div className="relative flex items-center justify-center">
        <div className="absolute top-1/2 left-2 -translate-y-1/2">
          <SearchIcon size={18} color="grey" />
        </div>
        <Input
          placeholder="주소 이동"
          value={searchInput.value}
          onChange={(e) => {
            searchInput.handleChange(e);

            if (e.target.value.length > 0) setResultModal(true);
            else setResultModal(false);
          }}
          onFocus={(e) => {
            if (e.target.value.length > 0) setResultModal(true);
            else setResultModal(false);
          }}
          onBlur={() => setResultModal(false)}
          className="rounded-sm border border-solid placeholder:text-grey-dark pl-8 text-base"
        />
        <button
          className="absolute top-1/2 right-2 -translate-y-1/2"
          onClick={searchInput.resetValue}
        >
          <ImCancelCircle />
        </button>
      </div>

      {resultModal && searchInput.value.length > 0 && (
        <div className="absolute top-10 left-0 w-full z-10 bg-black rounded-sm border border-solid p-4">
          {isError && <div>잠시 후 다시 시도해 주세요.</div>}
          {data?.documents.length === 0 && <div>검색 결과가 없습니다.</div>}
          {data?.documents.map((document) => {
            return (
              <MapSearchResult
                key={document.id}
                title={document.place_name}
                subTitle={document.address_name}
                lat={Number(document.y)}
                lng={Number(document.x)}
                reset={searchInput.resetValue}
              />
            );
          })}
        </div>
      )}
    </div>
  );
};

export default MapSearch;
