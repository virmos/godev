import useSWR from "swr"
import axios from "axios";

export const handler = () => () => {
    const swrRes = useSWR(`/api/products`,
        async () => {
            const { data } = await axios.get(`/api/products`);
            return data
        }
    )
  
    return swrRes
}