"use client";

import SignInBtn from "../account/SignInBtn";
import MemeInput from "./MemeInput";
import MemeList from "./list/MemeList";
import { useLatestMeme } from "@/hooks/galactica/meme/useMeme";
import { useCallback, useState, ReactElement } from "react";
import { OptimisticMemeUpload } from "@/api/galactica/meme/meme.type";
import MemeListItem from "./list/MemeListItem";
import { MemeFilter } from "./meme.type";

interface MemeTimelineProps {
  filter?: MemeFilter;
}

// Default to load latest approved memes
const MemeTimeline = ({
  filter = {
    status: "PENDING",
  },
}: MemeTimelineProps): ReactElement<MemeTimelineProps> => {
  const { memes, loadMore, isLoading, hasNext } = useLatestMeme(filter);

  const [userMemes, setUserMemes] = useState<OptimisticMemeUpload[]>([]);

  const handleLoadMore = async () => {
    if (
      isLoading ||
      !hasNext ||
      window.innerHeight + document.documentElement.scrollTop <
        document.documentElement.offsetHeight * 0.9
    ) {
      return;
    }

    await loadMore();
  };

  // Optimistically update the meme list UI after user upload their meme
  const handleMemeUpload = useCallback((meme: OptimisticMemeUpload) => {
    setUserMemes((userMemes) => {
      if (userMemes.length === 0) {
        return [meme];
      }

      return [...userMemes, meme];
    });
  }, []);

  return (
    <>
      <div className="flex items-center justify-between">
        <h1 className="text-4xl font-bold">#Truthmemes</h1>
        <SignInBtn />
      </div>
      <div className="mt-4">
        <MemeInput onUpload={handleMemeUpload} />
      </div>
      <div>
        {userMemes.map((meme, i) => (
          <MemeListItem {...meme} isServerMeme={false} key={i} />
        ))}
        <MemeList memes={memes} loadMore={handleLoadMore} />
        {isLoading && (
          <div className="py-4 flex items-center justify-center">
            Loading more...
          </div>
        )}
      </div>
    </>
  );
};

export default MemeTimeline;
