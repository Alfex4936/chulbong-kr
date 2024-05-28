"use client";

import SearchIcon from "@/components/icons/SearchIcon";
import { Input } from "@/components/ui/input";
import useInput from "@/hooks/common/useInput";
import useSearch from "@/hooks/query/search/useSearch";
import { useEffect, useRef, useState, type FocusEvent } from "react";
import { ImCancelCircle } from "react-icons/im";
import SearchResult from "./SearchResult";

interface Props {
  mini?: boolean;
  searchToggle?: boolean;
  sticky?: boolean;
}

const SearchInput = ({ sticky = false }: Props) => {
  const [query, setQuery] = useState("");
  const searchInput = useInput("");

  const { data: searchData, isError: searchError } = useSearch(query);

  const [resultModal, setResultModal] = useState(false);

  const searchResultRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handler = setTimeout(() => setQuery(searchInput.value), 300);

    return () => clearTimeout(handler);
  }, [searchInput.value]);

  const handleBlur = (e: FocusEvent<HTMLInputElement>) => {
    if (
      searchResultRef.current &&
      searchResultRef.current.contains(e.relatedTarget as Node)
    ) {
      return;
    }
    setResultModal(false);
  };

  return (
    <div
      className={`${
        sticky ? "sticky" : "relative"
      } top-0 left-0 mx-auto mb-4 bg-black z-[1000]`}
    >
      <div className="relative flex items-center justify-center">
        <div className="absolute top-1/2 left-2 -translate-y-1/2">
          <SearchIcon size={18} color="grey" />
        </div>
        <Input
          placeholder="장소, 위치를 입력하세요"
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
          // onBlur={handleBlur}
          className="rounded-sm border border-solid border-grey placeholder:text-grey-dark pl-8 text-base"
        />
        <button
          className="absolute top-1/2 right-2 -translate-y-1/2"
          onClick={searchInput.resetValue}
        >
          <ImCancelCircle />
        </button>
      </div>

      {resultModal && searchData && searchData.markers.length > 0 && (
        <div
          className="absolute top-10 left-0 w-full max-h-72 overflow-auto z-10 bg-black rounded-sm border border-solid 
                    scrollbar-thin border-grey p-4"
          ref={searchResultRef}
        >
          {searchError && <div>잠시 후 다시 시도해 주세요.</div>}
          {searchData.markers.length === 0 && <div>검색 결과가 없습니다.</div>}
          {searchData.markers.map((data) => {
            return (
              <SearchResult
                key={data.markerId}
                title={data.address}
                markerId={data.markerId}
              />
            );
          })}
        </div>
      )}
    </div>
  );
};

export default SearchInput;
