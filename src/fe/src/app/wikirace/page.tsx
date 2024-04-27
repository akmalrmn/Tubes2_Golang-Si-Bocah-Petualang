import React from "react";
import InputForm from "./components/input-form";
import Result from "./components/result";

function WikiRace() {
  return (
    <main className='flex flex-col justify-center items-center'>
      <div className=' py-24 px-12 z-[2] h-screen w-screen flex flex-col items-center space-y-4'>
        <div className='font-schoolbell bg-white w-fit relative text-center'>
          <div className='absolute -inset-16 bg-white blur-lg z-[-1] rounded-full'></div>
          <h1 className='text-[100px]'>WikiRace</h1>
          <p className='text-xl'>Find shortest path from</p>
        </div>

        <div className='relative bg-white p-8 font-schoolbell'>
          <div className='absolute -inset-16 bg-white blur-lg z-[-1] rounded-full'></div>
          <InputForm />
        </div>
      </div>

      <Result id='result' />
    </main>
  );
}

export default WikiRace;
