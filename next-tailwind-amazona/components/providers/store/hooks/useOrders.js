import useSWR from "swr"
import axios from "axios";

export const handler = () => () => {
    const swrRes = useSWR(`/api/orders`,
        async () => {
            const { data } = await axios.get(`/api/orders`);
            return data
        }
    )
  
    return swrRes
}