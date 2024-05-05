"use client";

import GrowBox from "@/components/atom/GrowBox";
import MinusIcon from "@/components/icons/MinusIcon";
import PlusIcon from "@/components/icons/PlusIcon";
import { useState } from "react";

interface FacilityProps {
  name: string;
  count: number;
  increase: VoidFunction;
  decrease: VoidFunction;
}

const FacilityList = ({ count, name, decrease, increase }: FacilityProps) => {
  return (
    <div className="flex items-center mb-2">
      <span>{name}</span>
      <GrowBox />
      <span className="flex items-center">
        <button
          className="rounded-full p-1 hover:bg-white-tp-dark"
          onClick={() => decrease()}
        >
          <MinusIcon size={18} />
        </button>
        <span className="mx-3">{count}</span>
        <button
          className="rounded-full p-1 hover:bg-white-tp-dark"
          onClick={() => increase()}
        >
          <PlusIcon size={18} />
        </button>
      </span>
    </div>
  );
};

const Facilities = () => {
  const [chulbong, setChulbong] = useState(0);
  const [penghang, setPenghang] = useState(0);

  return (
    <div>
      <FacilityList
        name="철봉"
        count={chulbong}
        increase={() => {
          if (chulbong === 99) return;
          setChulbong((prev) => prev + 1);
        }}
        decrease={() => {
          if (chulbong === 0) return;
          setChulbong((prev) => prev - 1);
        }}
      />
      <FacilityList
        name="평행봉"
        count={penghang}
        increase={() => {
          if (penghang === 99) return;
          setPenghang((prev) => prev + 1);
        }}
        decrease={() => {
          if (penghang === 0) return;
          setPenghang((prev) => prev - 1);
        }}
      />
    </div>
  );
};

export default Facilities;
