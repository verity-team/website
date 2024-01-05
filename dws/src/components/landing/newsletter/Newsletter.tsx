"use client";

import { subscribeEmail } from "@/api/dws/email/email";
import { UsersRound } from "lucide-react";
import { FormEvent, useRef, useState } from "react";
import toast from "react-hot-toast";

const Newsletter = () => {
  const [subscribed, setSubscribed] = useState(false);
  const emailRef = useRef<HTMLInputElement>(null);

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const email = emailRef.current?.value;
    if (email == null || email.trim() === "") {
      toast.error("Please input email before subscribe");
      return;
    }

    try {
      await subscribeEmail(email);
      setSubscribed(true);
    } catch (error: any) {
      toast.error(error.message);
    }
  };

  return (
    <div className="bg-cyellow mx-2">
      <div className="bg-corange py-8 flex flex-col items-center justify-center">
        <div>
          <h2 className="text-2xl inline-block break-words">
            Get notified. Get <span className="text-cred">Memes.</span>
          </h2>
        </div>
        <div>
          <p className="text-sm inline-block break-words font-sans font-bold text-center mt-2">
            No SPAM, just Memes and News
            <br />
            Unsubcribe any time.
          </p>
        </div>
        <form
          onSubmit={handleSubmit}
          className="w-full flex flex-col items-center px-4 mt-4"
        >
          <input
            className="p-4 w-full border-2 border-gray-400"
            placeholder="Enter your email"
            required
            type="email"
            ref={emailRef}
          />
          <button
            type="submit"
            className="mt-4 px-6 py-2 bg-cred text-white rounded-3xl border-4 border-black text-2xl tracking-wide uppercase"
          >
            <div className="flex items-center justify-center space-x-2">
              Submit
            </div>
          </button>
        </form>
      </div>
    </div>
  );
};

export default Newsletter;
