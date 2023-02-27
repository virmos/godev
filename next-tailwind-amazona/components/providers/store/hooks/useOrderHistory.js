import useSWR from "swr"
import axios from "axios";

export const handler = () => () => {
    const swrRes = useSWR(`/api/orders/history`,
        async () => {
            const { data } = await axios.get(`/api/orders/history`);
            return data
        }
    )
  
    return swrRes
}