"use client";

import useBookmarkData from "@/hooks/query/user/useBookmarkData";
import BookmarkList from "../_component/BookmarkList";

const BookmarkClient = () => {
  const { data: bookmarks } = useBookmarkData();

  return (
    <div>
      <div className="text-red text-center text-sm mb-4">
        위치는 총 10개까지 저장이 가능합니다.
      </div>
      <ul>
        {bookmarks?.map((bookmark) => {
          return (
            <BookmarkList
              key={bookmark.markerId}
              title={bookmark.address || "지원되지 않는 주소입니다."}
              subTitle={bookmark.description || ""}
              lng={bookmark.longitude}
              lat={bookmark.latitude}
              markerId={bookmark.markerId}
            />
          );
        })}
      </ul>
    </div>
  );
};

export default BookmarkClient;
