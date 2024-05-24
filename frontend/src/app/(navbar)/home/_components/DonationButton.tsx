import Link from "next/link";

const DonationButton = () => {
  return (
    <Link
      href={"https://toss.me/pullupko"}
      className="block w-full text-center bg-zinc-700 p-2 text-lg rounded-sm hover:bg-grey-dark-1"
      target="_black"
    >
      💸 후원하기
    </Link>
  );
};

export default DonationButton;
