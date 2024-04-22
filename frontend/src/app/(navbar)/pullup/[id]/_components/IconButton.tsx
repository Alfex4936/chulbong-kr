type Props = {
  icon: React.ReactNode;
  top?: number;
  left?: number;
  right?: number;
  bottom?: number;
};

const IconButton = ({ icon, top, left, right, bottom }: Props) => {
  return (
    <button
      className="absolute p-1 z-20 rounded-sm bg-white-tp-light hover:bg-white-tp-dark"
      style={{ top, left, right, bottom }}
    >
      {icon}
    </button>
  );
};

export default IconButton;
