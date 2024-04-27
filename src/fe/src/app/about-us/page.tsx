import PeopleCard from "./components/people-card";
import maulvi from "@/assets/images/maulvi.png";
import akmal from "@/assets/images/zayn.jpg";
import ojan from "@/assets/images/ojan.jpg";

const PeopleData = [
  {
    name: "Maulvi Ziadinda Maulana",
    nim: "13522122",
    desc: "I'm Maulvi, a passionate 4rd-semester undergraduate student pursuing Informatics Engineering at ITB. I am deeply engrossed in the world of front-end development and continuously honing my skills in this area to create engaging digital experiences. Currently, I'm expanding my skills into full-stack development, exploring the intricacies of back-end programming and databases. Check out my portofolio at maulvi-zm.github.io",
    image: maulvi,
  },
  {
    name: "Muhammad Fauzan Azhim",
    nim: "13522153",
    desc: "Hello! I'm Fauzan, a passionate Game Developer with a keen interest in Capture The Flag (CTF) competitions. Currently i'm in 4rd semester of Informatics Engineering, and i deeply involved in a mobile game project, which is yet to be named. If you want to contact me, sent dm at my instagram @fauzannazz",
    image: ojan,
  },
  {
    name: "Muhammad Akmal Ramadhan",
    nim: "13522161",
    desc: "Halo, saya akmal, saya ganteng.❤️",
    image: akmal,
  },
];

function AboutUs() {
  return (
      <main className='h-screen w-screen p-24'>
        <div className='flex flex-row w-full h-full justify-center items-center z-10'>
          {PeopleData.map((person, index) => (
              <PeopleCard key={index} {...person} />
          ))}
        </div>
      </main>
  );
}

export default AboutUs;
