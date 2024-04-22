type Props = { size?: number; color?: "black" | "white" };

const MinusIcon = ({ size = 24, color = "white" }: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      stroke={color === "black" ? "#222222" : "#F0F0F0"}
      fill="none"
    >
      <path
        d="M20 12H4"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
};

export default MinusIcon;
