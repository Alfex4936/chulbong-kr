import GrowBox from "@/components/atom/GrowBox";
import { LocationIcon } from "@/components/icons/LocationIcons";
import usePageLoadingStore from "@/store/usePageLoadingStore";
import { useRouter } from "next/navigation";

interface Props {
  title: string;
  markerId: number;
}

const SearchResult = ({ title, markerId }: Props) => {
  const router = useRouter();
  const { setLoading } = usePageLoadingStore();

  return (
    <button
      className={`flex w-full items-center p-4 rounded-sm mb-2 duration-100 hover:bg-zinc-700 hover:scale-95 last:mb-0`}
      onClick={() => {
        router.push(`/pullup/${markerId}`);
        setLoading(true);
      }}
    >
      <div className="w-3/4">
        <div className={`truncate text-left mr-2 text-sm`}>{title}</div>
      </div>
      <GrowBox />
      <div>
        <LocationIcon selected={false} size={18} />
      </div>
    </button>
  );
};

export default SearchResult;
