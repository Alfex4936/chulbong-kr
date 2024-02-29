import * as Styled from "./AddChinupSkeleton.style";

const AddChinupSkeleton = () => {
  return (
    <div style={{ display: "flex", flexDirection: "column" }}>
      <Styled.TitleSkeleton />
      <Styled.UploadSkeleton />
      <Styled.InputSkeleton />
      <Styled.ButtonSkeleton />
    </div>
  );
};

export default AddChinupSkeleton;
