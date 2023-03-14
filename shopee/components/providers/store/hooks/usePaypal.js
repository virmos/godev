import useSWR from "swr"
import axios from "axios";

export const handler = () => () => {
    const swrRes = useSWR(`/api/keys/paypal`,
        async () => {
            const { data } = await axios.get(`/api/keys/paypal`);
            return data
        }
    )
  
    return swrRes
}