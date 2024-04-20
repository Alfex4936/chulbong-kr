type Props = { selected: boolean; size?: number };

const ChatBubbleIcon = ({ selected, size = 35 }: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 40 40"
      fill="none"
    >
      <path
        d="M36.6666 19.2779C36.6666 28.0832 29.2036 35.2224 19.9999 35.2224C18.9178 35.2239 17.8386 35.1237 16.7756 34.9242C16.0105 34.7804 15.6279 34.7085 15.3608 34.7494C15.0937 34.7902 14.7152 34.9914 13.9582 35.394C11.8168 36.5329 9.31984 36.935 6.91844 36.4884C7.83115 35.3657 8.4545 34.0187 8.72955 32.5747C8.89622 31.6914 8.48325 30.8334 7.86474 30.2052C5.05547 27.3525 3.33325 23.5085 3.33325 19.2779C3.33325 10.4727 10.7962 3.33337 19.9999 3.33337C29.2036 3.33337 36.6666 10.4727 36.6666 19.2779Z"
        stroke={"#F0F0F0"}
        fill={selected ? "#F0F0F0" : "transparent"}
        strokeWidth="2"
        strokeLinejoin="round"
      />
      <path
        d="M19.9924 20H20.0074M26.6516 20H26.6666M13.3333 20H13.3482"
        stroke={selected ? "#222222" : "#F0F0F0"}
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
};

export default ChatBubbleIcon;
