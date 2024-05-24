const MyLocateOverlay = () => {
  return (
    <div className="w-4 h-4 bg-rose-500 border border-red rounded-full transform -translate-x-1/2 -translate-y-1/2">
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-0 w-0 border border-rose-500 border-solid rounded-full opacity-0 animate-ripple " />
    </div>
  );
};

export default MyLocateOverlay;
