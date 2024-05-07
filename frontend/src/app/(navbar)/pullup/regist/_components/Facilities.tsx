"use client";

import GrowBox from "@/components/atom/GrowBox";
import MinusIcon from "@/components/icons/MinusIcon";
import PlusIcon from "@/components/icons/PlusIcon";
import useUploadFormDataStore from "@/store/useUploadFormDataStore";
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
  const {
    facilities,
    increaseChulbong,
    decreaseChulbong,
    increasePenghang,
    decreasePenghang,
  } = useUploadFormDataStore();
  // const [chulbong, setChulbong] = useState(0);
  // const [penghang, setPenghang] = useState(0);

  return (
    <div>
      <FacilityList
        name="철봉"
        count={facilities.철봉}
        increase={() => {
          if (facilities.철봉 === 99) return;
          increaseChulbong();
          // setChulbong((prev) => prev + 1);
        }}
        decrease={() => {
          if (facilities.철봉 === 0) return;
          decreaseChulbong();
          // setChulbong((prev) => prev - 1);
        }}
      />
      <FacilityList
        name="평행봉"
        count={facilities.평행봉}
        increase={() => {
          if (facilities.평행봉 === 99) return;
          increasePenghang();
          // setPenghang((prev) => prev + 1);
        }}
        decrease={() => {
          if (facilities.평행봉 === 0) return;
          // setPenghang((prev) => prev - 1);
          decreasePenghang();
        }}
      />
    </div>
  );
};

export default Facilities;
