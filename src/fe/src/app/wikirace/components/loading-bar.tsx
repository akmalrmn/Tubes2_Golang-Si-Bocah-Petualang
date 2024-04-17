import React from "react";

function LoadingBar() {
  console.log("LoadingBar");
  const backgroundStyle = {
    backgroundImage:
      "linear-gradient(-45deg, #000010 25%, #666666 25%, #666666 50%, #000010 50%,#000010 75%, #666666 75%)",
    backgroundSize: "100px 100px",
  };
  return (
    <div
      className='flex justify-center items-center z-10 w-[300px] h-6 rounded-full ring-8 ring-slate-300 animate-sliding'
      style={backgroundStyle}
    ></div>
  );
}

export default LoadingBar;
