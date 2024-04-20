type Props = { selected: boolean; size?: number };

const NotificationIcon = ({ selected, size = 35 }: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 40 40"
      fill="none"
    >
      <path
        d="M8.59728 19.1517C8.47482 21.4784 8.6156 23.955 6.53688 25.514C5.5694 26.2397 5 27.3784 5 28.5879C5 30.2514 6.303 31.6667 8 31.6667H32C33.697 31.6667 35 30.2514 35 28.5879C35 27.3784 34.4307 26.2397 33.4632 25.514C31.3843 23.955 31.5252 21.4784 31.4027 19.1517C31.0835 13.0871 26.073 8.33337 20 8.33337C13.9269 8.33337 8.91647 13.0871 8.59728 19.1517Z"
        stroke={"#F0F0F0"}
        fill={selected ? "#F0F0F0" : "#222222"}
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M17.5 5.20837C17.5 6.58909 18.6193 8.33337 20 8.33337C21.3807 8.33337 22.5 6.58909 22.5 5.20837C22.5 3.82766 21.3807 3.33337 20 3.33337C18.6193 3.33337 17.5 3.82766 17.5 5.20837Z"
        stroke={"#F0F0F0"}
        fill={selected ? "#F0F0F0" : "#222222"}
        strokeWidth="2"
      />
      <path
        d="M25 31.6666C25 34.4281 22.7615 36.6666 20 36.6666C17.2385 36.6666 15 34.4281 15 31.6666"
        stroke={"#F0F0F0"}
        fill={selected ? "#F0F0F0" : "#222222"}
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
};

export default NotificationIcon;
