import { Maybe } from "@/utils/types";
import { CustomError, UserDonationData } from ".";

export type AffiliateCode = string | "none";

// Interface to construct request body for /wallet/connection
export interface WalletAffliateRequest {
  code: AffiliateCode;
  address: string;
}

export interface WalletAffiliateResponse {
  data: Maybe<UserDonationData>;
  error: Maybe<CustomError>;
}

export interface GenAffiliateRequest {
  address: string;
  timestamp: string;
  signature: string;
}

// Interface for /affliate/code response, return generated code to user
export interface GenAffliateResponse {
  address: string;
  code: AffiliateCode;
  ts: string;
}
