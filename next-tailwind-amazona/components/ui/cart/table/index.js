import { XCircleIcon } from "@heroicons/react/outline";
import Image from "next/image";
import Link from "next/link";

export default function CartTable({cartItems, updateCartHandler, removeItemHandler}) {
    return (
        <table className="min-w-full ">
            <thead className="border-b">
                <tr>
                    <th className="p-5 text-left">Item</th>
                    <th className="p-5 text-right">Quantity</th>
                    <th className="p-5 text-right">Price</th>
                    <th className="p-5">Action</th>
                </tr>
            </thead>
            <tbody>
                {cartItems.map((item) => (
                    <tr key={item.slug} className="border-b">
                        <td>
                            <Link href={`/product/${item.slug}`}>
                                <a className="flex items-center">
                                    <Image
                                        src={item.image}
                                        alt={item.name}
                                        width={50}
                                        height={50}
                                    ></Image>
                                    &nbsp;
                                    {item.name}
                                </a>
                            </Link>
                        </td>
                        <td className="p-5 text-right">
                            <select
                                value={item.quantity}
                                onChange={(e) =>
                                    updateCartHandler(item, e.target.value)
                                }
                            >
                                {[...Array(item.countInStock).keys()].map((x) => (
                                    <option key={x + 1} value={x + 1}>
                                        {x + 1}
                                    </option>
                                ))}
                            </select>
                        </td>
                        <td className="p-5 text-right">${item.price}</td>
                        <td className="p-5 text-center">
                            <button onClick={() => removeItemHandler(item)}>
                                <XCircleIcon className="h-5 w-5"></XCircleIcon>
                            </button>
                        </td>
                    </tr>
                ))}
            </tbody>
        </table>
    );
}
