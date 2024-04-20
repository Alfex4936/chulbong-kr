type Props = { selected: boolean; size?: number };

const UserCircleIcon = ({ selected, size = 35 }: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 40 40"
      fill="none"
    >
      <path
        d="M19.9999 36.6667C29.2047 36.6667 36.6666 29.2048 36.6666 20C36.6666 10.7953 29.2047 3.33337 19.9999 3.33337C10.7952 3.33337 3.33325 10.7953 3.33325 20C3.33325 29.2048 10.7952 36.6667 19.9999 36.6667Z"
        stroke={"#F0F0F0"}
        fill={selected ? "#F0F0F0" : "transparent"}
        strokeWidth="2"
      />
      <path
        d="M12.5 28.3333C16.3862 24.263 23.572 24.0713 27.5 28.3333M24.1585 15.8333C24.1585 18.1345 22.2903 20 19.9858 20C17.6815 20 15.8133 18.1345 15.8133 15.8333C15.8133 13.5321 17.6815 11.6666 19.9858 11.6666C22.2903 11.6666 24.1585 13.5321 24.1585 15.8333Z"
        stroke={selected ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
      />
    </svg>
  );
};

export default UserCircleIcon;
