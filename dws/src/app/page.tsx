import ClientRoot from "@/components/ClientRoot";
import Banner from "@/components/landing/banner/Banner";
import BannerSection from "@/components/landing/banner/BannerSection";
import MemeSlideshow from "@/components/landing/carousel/MemeSlideshow";
import LandingFooter from "@/components/landing/footer/LandingFooter";
import Navbar from "@/components/landing/navbar/Navbar";
import Newsletter from "@/components/landing/newsletter/Newsletter";
import Roadmap from "@/components/landing/roadmap/Roadmap";
import DonateFormV2 from "@/components/metamask/donatev2/DonateFormv2";
import dynamic from "next/dynamic";

const LaunchTimer = dynamic(
  () => import("@/components/landing/banner/LaunchTimer"),
  { ssr: false }
);

export default function Home() {
  return (
    <>
      <Navbar />
      <div className="mx-8">
        <div className="px-24 pt-10">
          <LaunchTimer />
        </div>
      </div>
      <Banner />
      <MemeSlideshow />
      <Roadmap />
      <BannerSection className="bg-white">
        <h1 className="text-6xl leading-loose-2xl italic">
          Uniting and tokenizing the truth movement.
        </h1>
        <h1 className="text-4xl leading-loose-xl my-6 break-words tracking-wide">
          Bringing together the two most vibrant and viral communities around
          truth-seeking: Meme-Artists and crypto-lovers.
        </h1>
      </BannerSection>
      <ClientRoot>
        <div className="flex items-center justify-center my-2">
          <DonateFormV2 />
        </div>
      </ClientRoot>
      <Newsletter />
      <LandingFooter />
    </>
  );
}
