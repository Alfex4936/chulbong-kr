type Props = { size?: number; color?: "black" | "white" };

const EditIcon = ({ size = 24, color = "white" }: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 15 16"
      fill="none"
    >
      <path
        d="M8.79606 2.92841C9.26181 2.4238 9.49469 2.1715 9.74213 2.02433C10.3392 1.66922 11.0744 1.65818 11.6815 1.9952C11.9331 2.13487 12.1731 2.38007 12.6531 2.87047C13.1332 3.36087 13.3733 3.60607 13.5099 3.86308C13.8399 4.4832 13.8291 5.23424 13.4814 5.84419C13.3374 6.09698 13.0904 6.33486 12.5964 6.81063L6.71913 12.4714C5.78305 13.3731 5.315 13.8239 4.73004 14.0523C4.14508 14.2808 3.502 14.264 2.21585 14.2304L2.04086 14.2258C1.64932 14.2156 1.45354 14.2104 1.33974 14.0812C1.22594 13.9521 1.24148 13.7527 1.27255 13.3539L1.28943 13.1373C1.37688 12.0147 1.42061 11.4534 1.63982 10.9489C1.85903 10.4443 2.23715 10.0347 2.99339 9.21531L8.79606 2.92841Z"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinejoin="round"
      />
      <path
        d="M8.125 3L12.5 7.375"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinejoin="round"
      />
      <path
        d="M8.75 14.25H13.75"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
};

export default EditIcon;
