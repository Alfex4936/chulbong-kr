import CircularProgress from "@mui/material/CircularProgress";
import CenterBox from "../CenterBox/CenterBox";

const Loader = () => {
  return (
    <CenterBox bg="black">
      <CircularProgress size={50} sx={{ color: "#fff" }} />
    </CenterBox>
  );
};

export default Loader;
