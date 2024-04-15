type Props = { selected: boolean; size?: number };

const HomeIcon = ({ selected, size = 35 }: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 37 36"
      fill="none"
    >
      <path
        d="M12.9572 34.3333L12.5291 28.4851C12.2989 25.341 14.8514 22.6666 18.0822 22.6666C21.3131 22.6666 23.8656 25.341 23.6355 28.4851L23.2072 34.3333"
        stroke={selected ? "#222222" : "#F0F0F0"}
        strokeWidth="2"
      />
      <path
        d="M1.60029 19.6892C0.997216 15.8603 0.695678 13.946 1.43762 12.249C2.17955 10.5519 3.82564 9.39072 7.11778 7.06843L9.57752 5.33333C13.6729 2.44445 15.7207 1 18.0833 1C20.4461 1 22.4937 2.44445 26.5891 5.33333L29.0489 7.06843C32.3411 9.39072 33.9872 10.5519 34.7291 12.249C35.4711 13.946 35.1695 15.8603 34.5663 19.6892L34.0521 22.954C33.1973 28.3815 32.7697 31.0953 30.7745 32.7143C28.7793 34.3333 25.8625 34.3333 20.0287 34.3333H16.1379C10.3042 34.3333 7.3873 34.3333 5.39212 32.7143C3.39695 31.0953 2.96949 28.3815 2.11457 22.954L1.60029 19.6892Z"
        stroke={selected ? "#222222" : "#F0F0F0"}
        strokeWidth="2"
        strokeLinejoin="round"
      />
    </svg>
  );
};
export default HomeIcon;
