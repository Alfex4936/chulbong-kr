"use client";

import MarkerListItem from "@/components/atom/MarkerListItem";
import SearchIcon from "@/components/icons/SearchIcon";
import { Input } from "@/components/ui/input";
import useInput from "@/hooks/common/useInput";
import useSearchLocationData from "@/hooks/query/useSearchLocationData";
import { useEffect, useState } from "react";
// TODO: input fontSize 16px 이상 https://velog.io/@jungsu/iOSSafari-input-%EC%9D%98-focus-%ED%99%95%EB%8C%80-%EB%B0%A9%EC%A7%80%ED%95%98%EA%B8%B0

const SearchInput = () => {
  const [query, setQuery] = useState("");
  const searchInput = useInput("");

  const { data, isError } = useSearchLocationData(query);

  useEffect(() => {
    const handler = setTimeout(() => setQuery(searchInput.value), 300);

    return () => clearTimeout(handler);
  }, [searchInput.value]);

  return (
    <div className="relative mx-auto mb-4">
      <div className="relative flex items-center justify-center">
        <div className="absolute top-1/2 left-2 -translate-y-1/2">
          <SearchIcon size={18} color="grey" />
        </div>
        <Input
          placeholder="장소, 위치를 입력하세요"
          value={searchInput.value}
          onChange={searchInput.handleChange}
          className="rounded-sm border border-solid border-grey placeholder:text-grey-dark pl-8 text-base"
        />
      </div>
      {searchInput.value.length > 0 && (
        <div className="absolute top-10 left-0 w-full z-10 bg-black rounded-sm border border-solid border-grey p-4">
          {isError && <div>잠시 후 다시 시도해 주세요.</div>}
          {data?.documents.length === 0 && <div>검색 결과가 없습니다.</div>}
          {data?.documents.map((document) => {
            return (
              <MarkerListItem
                key={document.id}
                title={document.place_name}
                subTitle={document.address_name}
                lat={Number(document.y)}
                lng={Number(document.x)}
              />
            );
          })}
        </div>
      )}
    </div>
  );
};

export default SearchInput;
