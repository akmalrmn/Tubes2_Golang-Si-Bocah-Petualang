import React from "react";
import { Card } from "@/components/card";
import Image from "next/image";
import { StaticImport } from "next/dist/shared/lib/get-img-props";

interface PeopleCardProps {
  image: StaticImport;
  name: string;
  nim: string;
  desc: string;
}

function PeopleCard({ image, name, nim, desc }: PeopleCardProps) {
  return (
    <Card
      icon={
        <div className='rounded-full'>
          <Image
            src={image}
            alt='logo'
            height={300}
            width={300}
            style={{
              maxHeight: "300px",
              maxWidth: "300px",
              objectFit: "cover",
            }}
            className='rounded-full'
          />
        </div>
      }
      text={
        <>
          <p>{name}</p>
          <p>{nim}</p>
          <p className='text-[16px] font-light mt-10'>{desc}</p>
        </>
      }
    >
      <div className='w-full h-full bg-black'></div>
    </Card>
  );
}

export default PeopleCard;
