import React from "react";
import RandomText from "./components/random-text";

export default function Loading(
  props: React.DetailedHTMLProps<
    React.HTMLAttributes<HTMLDivElement>,
    HTMLDivElement
  >
) {
  return (
    <div
      {...props}
      className='min-h-screen min-w-[100vw] flex flex-col items-center justify-center bg-white z-10 relative gap-8'
    >
      <LoadingBar />
      <RandomText />
      <div className='w-full h-full absolute -inset-2 bg-white blur-md'></div>
    </div>
  );
}

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
