type Props = { size?: number; color?: "black" | "white" };

const ExitIcon = ({ size = 24, color = "white" }: Props) => {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
    >
      <path
        d="M10.247 6.74032C11.0733 7.56669 11.4865 7.97988 12 7.97988C12.5135 7.97988 12.9267 7.5667 13.753 6.74034L15.5066 4.98681C15.9142 4.57923 16.1181 4.37536 16.3301 4.25295C17.3964 3.63729 18.2747 4.24833 19.0132 4.98681C19.7517 5.7253 20.3627 6.60358 19.7471 7.66993C19.6246 7.88195 19.4208 8.08575 19.0133 8.49334L17.2599 10.2467C16.4335 11.0731 16.0201 11.4865 16.0201 12C16.0201 12.5135 16.4333 12.9267 17.2597 13.7531L19.0132 15.5066C19.4208 15.9142 19.6246 16.118 19.7471 16.33C20.3627 17.3964 19.7517 18.2747 19.0132 19.0132C18.2748 19.7517 17.3964 20.3627 16.3301 19.747C16.118 19.6247 15.9142 19.4209 15.5066 19.0132L13.7533 17.2599C12.9272 16.4337 12.5135 16.0201 12 16.0201C11.4865 16.0201 11.0729 16.4337 10.2467 17.2599L8.49341 19.0132C8.08577 19.4209 7.88196 19.6247 7.66993 19.747C6.60365 20.3627 5.72522 19.7517 4.98681 19.0132C4.24827 18.2747 3.63732 17.3964 4.25295 16.33C4.37537 16.118 4.57918 15.9142 4.98681 15.5066L6.74032 13.7531C7.56669 12.9267 7.97987 12.5135 7.97987 12C7.97987 11.4865 7.56647 11.0731 6.7401 10.2467L4.98673 8.49334C4.57915 8.08575 4.37536 7.88195 4.25295 7.66993C3.63729 6.60358 4.24833 5.7253 4.98681 4.98681C5.7253 4.24833 6.60357 3.63729 7.66993 4.25295C7.88195 4.37536 8.08581 4.57922 8.4934 4.9868L10.247 6.74032Z"
        stroke={color === "black" ? "#222222" : "#F0F0F0"}
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
};

export default ExitIcon;
