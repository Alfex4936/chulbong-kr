type Props = { size?: number; color?: "black" | "white" };

export const GpsOnIcon = ({ size = 24, color = "white" }: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
    >
      <path
        d="M20.5137 12C20.5137 16.6944 16.7081 20.5 12.0137 20.5C7.31925 20.5 3.51367 16.6944 3.51367 12C3.51367 7.30558 7.31925 3.5 12.0137 3.5C16.7081 3.5 20.5137 7.30558 20.5137 12Z"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
      />
      <path
        d="M22.5 12H20.5"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M3.5 12H1.5"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M12 1.5V3.5"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M12 20.5V22.5"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M15 12C15 13.6569 13.6568 15 12 15C10.3431 15 9 13.6569 9 12C9 10.3431 10.3431 9 12 9C13.6568 9 15 10.3431 15 12Z"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
      />
    </svg>
  );
};
export const GpsOffIcon = ({ size = 24, color = "white" }: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
    >
      <path
        d="M20.5137 12C20.5137 16.6944 16.7081 20.5 12.0137 20.5C7.31925 20.5 3.51367 16.6944 3.51367 12C3.51367 7.30558 7.31925 3.5 12.0137 3.5C16.7081 3.5 20.5137 7.30558 20.5137 12Z"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
      />
      <path
        d="M15.0003 9L9.00024 15M15.0003 15L9.00024 9"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
      />
      <path
        d="M22.5 12H20.5"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M3.5 12H1.5"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M12 1.5V3.5"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M12 20.5V22.5"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
};
