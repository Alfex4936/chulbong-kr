const MapLoading = () => {
  return (
    <div className="h-full w-full bg-black bg-grid-white/[0.2] relative flex items-center justify-center">
      <div
        className="absolute pointer-events-none inset-0 flex items-center justify-center 
                bg-black [mask-image:radial-gradient(ellipse_at_center,transparent_20%,black)]"
      />
      <p
        className="text-2xl font-bold relative z-20 bg-clip-text text-transparent 
                    bg-gradient-to-b from-neutral-200 to-neutral-500 py-8"
      >
        지도 만드는 중...
      </p>
    </div>
  );
};

export default MapLoading;
