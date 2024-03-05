import { useState } from "react";
import * as Styled from "./AroundMarker.style";
import IconButton from "@mui/material/IconButton";
import Tooltip from "@mui/material/Tooltip";
import SearchIcon from "@mui/icons-material/Search";

const AroundMarker = () => {
  const [value, setValue] = useState(100);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValue(Number(event.target.value));
  };

  return (
    <div>
      <Styled.RangeContainer>
        <p style={{ fontSize: ".8rem" }}>주변 {value}m</p>
        <input
          type="range"
          min="100"
          max="500"
          step="100"
          value={value}
          onChange={handleChange}
        />
        <Tooltip title="검색" arrow disableInteractive>
          <IconButton
            onClick={() => {
              console.log(1);
            }}
            aria-label="delete"
            sx={{
              color: "#333",
              width: "30px",
              height: "30px",
            }}
          >
            <SearchIcon sx={{ fontSize: 22 }} />
          </IconButton>
        </Tooltip>
      </Styled.RangeContainer>
    </div>
  );
};

export default AroundMarker;
