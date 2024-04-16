type Props = { size?: number; color?: "black" | "white" | "grey" };

const SearchIcon = ({ size = 24, color = "white" }: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 26 25"
      fill="none"
    >
      <path
        d="M18.8735 18.2292L23.604 22.9167"
        stroke={
          color === "black"
            ? "#222222"
            : color === "white"
            ? "#F0F0F0"
            : "#bebebe"
        }
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M21.5015 11.4583C21.5015 6.28059 17.2657 2.08325 12.0405 2.08325C6.81534 2.08325 2.5795 6.28059 2.5795 11.4583C2.5795 16.636 6.81534 20.8333 12.0405 20.8333C17.2657 20.8333 21.5015 16.636 21.5015 11.4583Z"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="2"
        strokeLinejoin="round"
      />
    </svg>
  );
};

export default SearchIcon;
