type Props = {
  children: React.ReactNode;
  center?: boolean;
};

const BlackLightBox = ({ children, center = false }: Props) => {
  return (
    <div
      className={`bg-black-light-2 mx-auto w-[90%] p-4 rounded-md ${
        center ? "text-center" : "text-left"
      }`}
    >
      {children}
    </div>
  );
};

export default BlackLightBox;
