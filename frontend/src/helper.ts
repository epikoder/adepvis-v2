import { LogDebug } from "../wailsjs/runtime/runtime";

type NullError = Error | null;
type Value = any;

interface GoResponse {
  Value: Value;
  Err: NullError;
}

export const invoke = async <T extends Value = null, E extends Error = Error>(
  fn: () => Promise<T>
): Promise<GoResponse> => {
  const res: GoResponse = {
    Value: null,
    Err: null,
  };

  try {
    const v = await fn();
    res.Value = <T>v;
    return res;
  } catch (error) {
    res.Err = new Error(<string>error);
    return res;
  }
};
