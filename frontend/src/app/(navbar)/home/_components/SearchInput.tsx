"use client";

import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import useSearchLocationData from "@/hooks/query/useSearchLocationData";
import { useState } from "react";

// TODO: 검색 입력창 스타일링 하기

const SearchInput = () => {
  const [value, setValue] = useState("");

  const { data } = useSearchLocationData(value);

  const handleChange = (text: string) => {
    setValue(text);
  };

  return (
    <Command className="rounded-md border shadow-md mx-auto mt-8 w-5/6 bg-black">
      <CommandInput
        className="text-grey"
        placeholder="검색어를 입력하세요..."
        value={value}
        onValueChange={handleChange}
      />
      <CommandList>
        <CommandEmpty className="text-grey-dark p-2">
          찾으시는 결과가 없습니다.
        </CommandEmpty>
        <CommandGroup heading="추천">
          <CommandItem>
            <span>서울</span>
          </CommandItem>
          <CommandItem>
            <span>남산 타워</span>
          </CommandItem>
          <CommandItem>
            <span>수원</span>
          </CommandItem>
        </CommandGroup>
      </CommandList>
    </Command>
  );
};

export default SearchInput;
