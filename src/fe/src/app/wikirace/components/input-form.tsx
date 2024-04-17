"use client";

import React, { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { SubmitHandler, set, useForm } from "react-hook-form";
import TurnArrow from "./arrow-left-right";
import { revalidatePath, revalidateTag } from "next/cache";
import Checkmark from "./checkmark";

interface Inputs {
  start: string;
  target: string;
  algorithm: string;
}

function InputForm() {
  const router = useRouter();
  const params = useSearchParams();

  const { register, handleSubmit, setValue, getValues } = useForm<Inputs>();

  const onSubmit: SubmitHandler<Inputs> = (data) => {
    router.push(
      `/wikirace/?start=${data.start}&target=${data.target}&algorithm=${data.algorithm}`
    );
  };

  const onArrowClick = () => {
    console.log("clicked", getValues("start"), getValues("target"));
    const temp = getValues("start");
    setValue("start", getValues("target"));
    setValue("target", temp);
  };

  useEffect(() => {
    if (params.has("start") && params.has("target")) {
      setValue("start", params.get("start") as string);
      setValue("target", params.get("target") as string);
    }
  }, []);

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className='flex flex-col space-y-4 items-center text-xl'
    >
      <div className='flex flex-row space-x-8 items-center'>
        <input
          {...register("start")}
          placeholder='Start'
          className='p-2 border-2 border-black w-56'
        />

        <button
          type='button'
          className='w-fit h-fit hover:scale-110'
          onClick={onArrowClick}
        >
          <TurnArrow />
        </button>

        <input
          {...register("target")}
          placeholder='Target'
          className='p-2 border-2 border-black w-56'
        />
      </div>

      <div className='text-gray-700 text-xl pt-6'>
        <p>Choose Algorithm:</p>
      </div>

      <div className='flex gap-8 text-xl'>
        <div className='inline-flex items-center'>
          <label
            className='relative flex items-center p-3 rounded-full cursor-pointer'
            htmlFor='algorithm'
            id='algorithm'
          >
            <input
              {...register("algorithm")}
              type='radio'
              value={"bfs"}
              className="before:content[''] peer relative h-5 w-5 cursor-pointer appearance-none rounded-full border border-blue-gray-200 text-gray-900 transition-all before:absolute before:top-2/4 before:left-2/4 before:block before:h-12 before:w-12 before:-translate-y-2/4 before:-translate-x-2/4 before:rounded-full before:bg-blue-gray-500 before:opacity-0 before:transition-opacity checked:border-gray-900 checked:before:bg-gray-900 hover:before:opacity-10"
            />
            <span className='absolute text-gray-900 transition-opacity opacity-0 pointer-events-none top-2/4 left-2/4 -translate-y-2/4 -translate-x-2/4 peer-checked:opacity-100 scale-125'>
              <Checkmark />
            </span>
          </label>
          <label
            className='mt-px font-light text-gray-700 cursor-pointer select-none'
            htmlFor='algorithm'
            id='algorithm'
          >
            Breath First Search
          </label>
        </div>

        <div className='inline-flex items-center'>
          <label
            className='relative flex items-center p-3 rounded-full cursor-pointer'
            htmlFor='html'
          >
            <input
              {...register("algorithm")}
              value={"ids"}
              type='radio'
              className="before:content[''] peer relative h-5 w-5 cursor-pointer appearance-none rounded-full border border-blue-gray-200 text-gray-900 transition-all before:absolute before:top-2/4 before:left-2/4 before:block before:h-12 before:w-12 before:-translate-y-2/4 before:-translate-x-2/4 before:rounded-full before:bg-blue-gray-500 before:opacity-0 before:transition-opacity checked:border-gray-900 checked:before:bg-gray-900 hover:before:opacity-10"
            />
            <span className='absolute text-gray-900 transition-opacity opacity-0 pointer-events-none top-2/4 left-2/4 -translate-y-2/4 -translate-x-2/4 peer-checked:opacity-100 scale-125'>
              <Checkmark />
            </span>
          </label>
          <label
            className='mt-px font-light text-gray-700 cursor-pointer select-none'
            htmlFor='html'
          >
            Iterative Deepening Search
          </label>
        </div>
      </div>

      <button type='submit' className='p-2 bg-black text-white w-1/5'>
        Submit
      </button>
    </form>
  );
}

export default InputForm;
