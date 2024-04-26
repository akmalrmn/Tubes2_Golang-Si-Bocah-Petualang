"use client";

import React, { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { SubmitHandler, UseFormSetValue, useForm } from "react-hook-form";
import TurnArrow from "../../../components/icons/arrow-left-right";
import Checkmark from "../../../components/icons/checkmark";
import axios from "axios";
import useDebounce from "../../../hooks/useDebounce";

interface Inputs {
  start: string;
  target: string;
  algorithm: string;
}

export default function InputForm() {
  const router = useRouter();
  const params = useSearchParams();

  const { register, handleSubmit, setValue, getValues } = useForm<Inputs>({
    defaultValues: {
      start: "",
      target: "",
      algorithm: "bfs",
    },
  });

  const onSubmit: SubmitHandler<Inputs> = (data) => {
    router.push(
      `/wikirace/?start=${data.start}&target=${data.target}&algorithm=${data.algorithm}`
    );
  };

  const onArrowClick = () => {
    const temp = getValues("start");
    setValue("start", getValues("target"));
    setValue("target", temp);
  };

  useEffect(() => {
    if (params.has("start") && params.has("target")) {
      setValue("start", params.get("start") as string);
      setValue("target", params.get("target") as string);
    }

    if (params.has("algorithm") && params.get("algorithm") === "ids") {
      setValue("algorithm", "ids");
    } else {
      setValue("algorithm", "bfs");
    }
  }, []);

  return (
    <form
      onSubmit={handleSubmit(onSubmit)}
      className='flex flex-col space-y-4 items-center text-xl'
    >
      <div className='flex flex-row space-x-8 items-center'>
        <div className='w-fit relative'>
          <InputSuggestion
            register={{
              ...register("start", {
                required: true,
              }),
              placeholder: "Start",
            }}
            type='start'
            setValue={setValue}
          />
        </div>

        <button
          type='button'
          className='w-fit h-fit hover:scale-110'
          onClick={onArrowClick}
        >
          <TurnArrow />
        </button>

        <div className='w-fit relative'>
          <InputSuggestion
            register={{
              ...register("target", {
                required: true,
              }),
              placeholder: "Target",
            }}
            type='target'
            setValue={setValue}
          />
        </div>
      </div>

      <div className='text-gray-700 text-xl pt-6'>
        <p>Choose Algorithm:</p>
      </div>

      <div className='flex gap-8 text-xl'>
        <div className='inline-flex items-center'>
          <CheckBoxInput
            text='Breadth First Search'
            register={{
              ...register("algorithm"),
              value: "bfs",
              type: "radio",
            }}
          />
        </div>

        <div className='inline-flex items-center'>
          <CheckBoxInput
            text='Iterative Deepening Search'
            register={{
              ...register("algorithm"),
              value: "ids",
              type: "radio",
            }}
          />
        </div>
      </div>

      <button
        type='submit'
        className='p-2 bg-black text-white w-1/5 disabled:opacity-50 disabled:cursor-not-allowed'
      >
        Submit
      </button>
    </form>
  );
}

function CheckBoxInput({
  text,
  register,
}: {
  text: string;
  register: React.DetailedHTMLProps<
    React.InputHTMLAttributes<HTMLInputElement>,
    HTMLInputElement
  >;
}) {
  return (
    <>
      <label
        className='relative flex items-center p-3 rounded-full cursor-pointer'
        htmlFor='html'
      >
        <input
          {...register}
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
        {text}
      </label>
    </>
  );
}

function InputSuggestion({
  register,
  type,
  setValue,
}: {
  register: React.DetailedHTMLProps<
    React.InputHTMLAttributes<HTMLInputElement>,
    HTMLInputElement
  >;
  type: "start" | "target";
  setValue: UseFormSetValue<Inputs>;
}) {
  const [search, setSearch] = React.useState<string>("");
  const debouncedValue = useDebounce(search);
  const [data, setData] = React.useState<string[] | null>(null);

  function handleChange(event: React.ChangeEvent<HTMLInputElement>) {
    setSearch(event.target.value);
  }

  function handleClick(title: string) {
    setValue(type, title);

    setData(null);
  }

  async function fetchData() {
    const queryParams = {
      action: "query",
      format: "json",
      gpssearch: debouncedValue,
      generator: "prefixsearch",
      prop: "pageprops|pageimages|pageterms",
      redirects: "",
      ppprop: "displaytitle",
      piprop: "thumbnail",
      pithumbsize: "160",
      pilimit: "30",
      wbptterms: "description",
      gpsnamespace: 0,
      gpslimit: 5,
      origin: "*",
    };

    try {
      const response = await axios({
        method: "get",
        params: queryParams,
        url: process.env.WIKIPEDIA_API_URL,
        headers: {
          "Api-User-Agent": "Tubes 2 Stima ITB; 13522122@std.stei.itb.ac.id",
        },
      });

      const title: string[] = [];

      for (const key in response.data.query.pages) {
        title.push(response.data.query.pages[key].title);
      }

      return title;
    } catch (e) {
      console.error(e);
      return null;
    }
  }

  useEffect(() => {
    if (debouncedValue) {
      fetchData().then((data) => {
        if (data) {
          setData(data);
        }
      });
    }
  }, [debouncedValue]);

  return (
    <>
      <input
        {...register}
        onChange={handleChange}
        className='p-2 border-2 border-black w-56'
      />
      {data && (
        <div className='h-fit w-full bg-white absolute border-black border-2 border-t-0 z-30'>
          {data.map((title) => (
            <div
              key={title}
              className='p-2 text-black hover:bg-slate-400 cursor-pointer'
              onClick={() => handleClick(title)}
            >
              {title}
            </div>
          ))}
        </div>
      )}
    </>
  );
}
