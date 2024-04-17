import React from "react";
import Image from "next/image";

type Props = {};

const images = ["/11.png", "/12.png", "/13.png", "/14.png", "/15.png"];

const ImageList = (props: Props) => {
  return (
    <div className="flex">
      <div className="w-1/2">
        {images.map((img, i) => {
          if (i % 2 === 1) return;
          return (
            <button>
              <img src={img} className="w-full" key={i} />
            </button>
          );
        })}
      </div>
      <div className="w-1/2">
        {images.map((img, i) => {
          if (i % 2 !== 1) return;
          return (
            <button>
              <img src={img} className="w-full" key={i} />
            </button>
          );
        })}
      </div>
    </div>
  );
};

export default ImageList;
