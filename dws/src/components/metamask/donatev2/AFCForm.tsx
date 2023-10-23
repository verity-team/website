"use client";

import {
  getUserDonationData,
  requestNewAffiliateCode,
  useUserDonationData,
} from "@/utils/api/clientAPI";
import { requestSignature } from "@/utils/metamask/sign";
import { connectWallet } from "@/utils/metamask/wallet";
import { Maybe } from "@/utils/types";
import { getRFC3339String } from "@/utils/utils";
import { useSDK } from "@metamask/sdk-react";
import {
  Dialog,
  DialogTitle,
  DialogContent,
  CircularProgress,
} from "@mui/material";
import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import { ReactElement, memo, useEffect, useMemo, useState } from "react";
import toast from "react-hot-toast";

interface AFCFormProps {
  account: string;
}

const AFCForm = ({ account }: AFCFormProps): ReactElement<AFCFormProps> => {
  const [isFormOpen, setFormOpen] = useState(false);
  const [userCode, setUserCode] = useState("");

  const { sdk } = useSDK();
  const { data, error, isLoading } = useUserDonationData(account);

  useEffect(() => {
    if (data == null || error != null) {
      return;
    }

    setUserCode(data.stats.affliate_code);
  }, [data, error]);

  const sharableLink = useMemo(() => {
    if (typeof window === "undefined") {
      return "";
    }

    const baseUrl = window.location.href.split("?")[0];
    if (userCode == null || userCode.trim() === "") {
      return baseUrl;
    }

    return `${baseUrl}?afc=${userCode}`;
  }, [userCode]);

  const handleOpenForm = () => {
    setFormOpen(true);
  };

  const handleCloseForm = () => {
    setFormOpen(false);
  };

  const handleGenAffiliateCode = async () => {
    // Try to (re)connect when there are no connected accounts
    let currentAccount = account;
    if (!currentAccount) {
      const wallet = await connectWallet(sdk);
      if (wallet == null) {
        return;
      }
      currentAccount = wallet;

      const userDonationData = await getUserDonationData(currentAccount);

      if (userDonationData?.user_data.affiliate_code != null) {
        setUserCode(userDonationData.user_data.affiliate_code);
        return;
      }
    }

    const currentDate = new Date();

    // Timestamp in seconds
    const timestamp = Math.floor(currentDate.getTime() / 1000);
    const messageDate = getRFC3339String(currentDate);
    const message = `get affiliate code, ${messageDate}`;

    let signature: Maybe<string> = null;
    try {
      signature = await requestSignature(currentAccount, message);
      console.log(signature);
    } catch {
      toast.error("Transaction rejected");
      return;
    }

    if (signature == null) {
      return;
    }

    try {
      const affiliateResponse = await requestNewAffiliateCode({
        address: currentAccount,
        timestamp,
        signature,
      });
      if (affiliateResponse == null) {
        toast.error("Server fail to generate new affiliate code");
        return;
      }

      setUserCode(affiliateResponse.code);
    } catch {
      toast.error("Server fail to generate new affiliate code");
    }
  };

  const handleCopyToClipboard = () => {
    navigator.clipboard.writeText(sharableLink);
    toast.success("Copied to clipboard");
  };

  return (
    <>
      <div className="text-center text-4xl">
        Hype?{" "}
        <span
          className="underline cursor-pointer hover:text-blue-700"
          onClick={handleOpenForm}
        >
          Share now!
        </span>
      </div>
      <Dialog
        open={isFormOpen}
        onClose={handleCloseForm}
        fullWidth={true}
        maxWidth="sm"
        className="rounded-lg"
      >
        <DialogTitle className="text-2xl">Share</DialogTitle>
        <DialogContent className="font-sans w-full">
          <div>
            {!!userCode ? (
              <>
                <div className="text-lg mb-2">Your affiliate code:</div>
                <input
                  value={userCode}
                  disabled
                  className="px-4 py-2 border-2 border-black rounded-lg w-full text-center"
                />
              </>
            ) : (
              <>
                {isLoading ? (
                  <div className="w-full flex items-center justify-center py-2 border-2 border-black bg-blue-600 rounded-lg">
                    <CircularProgress size={24} className="text-white" />
                  </div>
                ) : (
                  <button
                    className="w-full px-4 py-2 border-2 border-black bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                    onClick={handleGenAffiliateCode}
                  >
                    Generate affiliate code
                  </button>
                )}
                <div className="text-base italic mt-2">
                  Note: You will need to sign a message for verification
                  purposes
                </div>
              </>
            )}

            <div className="mt-4">
              <div className="text-lg mb-2">Sharable link:</div>
              <div className="px-4 py-2 border-2 border-black rounded-lg w-full flex items-center justify-between">
                <input value={sharableLink} disabled className="w-full" />
                <div onClick={handleCopyToClipboard}>
                  <ContentCopyIcon className="text-gray-600 hover:text-gray-900 cursor-pointer" />
                </div>
              </div>
              {/* <div className="flex flex-col items-center justify-center">
                <QRCode value={sharableLink} size={100} />
              </div> */}
            </div>
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
};

export default memo(AFCForm);
