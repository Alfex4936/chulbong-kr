"use client";

import { globeArcs, globeConfig } from "@/data/globeConfig";
import { motion } from "framer-motion";
import dynamic from "next/dynamic";
import { useRouter } from "next/navigation";

const World = dynamic(() => import("../ui/globe").then((m) => m.World), {
  ssr: false,
});

export const GlobeWrap = () => {
  const router = useRouter();

  return (
    <div className="flex flex-row items-center justify-center h-screen relative w-full">
      <div className="max-w-7xl mx-auto w-full relative overflow-hidden md:h-[40rem] px-4">
        <motion.div
          initial={{
            opacity: 0,
            y: 20,
          }}
          animate={{
            opacity: 1,
            y: 0,
          }}
          transition={{
            duration: 2,
          }}
          className="div mb-10"
        >
          <h2 className="text-center text-xl md:text-4xl font-bold text-black dark:text-white">
            대한민국 모든 지도
          </h2>
          <p className="text-center text-base md:text-lg font-normal text-neutral-700 dark:text-neutral-200 max-w-md mt-2 mx-auto">
            다른 사람들과 함께 원하는 위치를 등록하고 공유하세요!
          </p>
        </motion.div>
        <motion.div
          initial={{
            opacity: 0,
            y: 20,
          }}
          animate={{
            opacity: 1,
            y: 0,
          }}
          transition={{
            duration: 3,
          }}
          className="div mb-5 text-center"
        >
          <button
            className="relative inline-flex h-12 overflow-hidden rounded-xl p-[1px] focus:outline-none focus:ring-2 
          focus:ring-slate-400 focus:ring-offset-2 focus:ring-offset-slate-50"
            onClick={() => router.push("/home")}
          >
            <span className="absolute inset-[-1000%] animate-[spin_2s_linear_infinite] bg-[conic-gradient(from_90deg_at_50%_50%,#E2CBFF_0%,#393BB2_50%,#E2CBFF_100%)]" />
            <span className="inline-flex h-full w-full cursor-pointer items-center justify-center rounded-xl bg-slate-950 px-4 py-0 text-sm font-medium text-white backdrop-blur-3xl">
              철봉 지도 바로 가기
            </span>
          </button>
        </motion.div>

        <div className="h-72 md:h-full pb-10 z-10 md:pb-56">
          <World data={globeArcs} globeConfig={globeConfig} />
        </div>
      </div>
    </div>
  );
};
