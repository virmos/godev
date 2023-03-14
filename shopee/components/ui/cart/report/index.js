import { useRouter } from "next/router";

export default function CartReport({cartItems}) {
    const router = useRouter();
    return (
        <div className="card p-5">
            <ul>
                <li>
                    <div className="pb-3 text-xl">
                        Subtotal (
                        {cartItems.reduce((a, c) => a + c.quantity, 0)}) : $
                        {cartItems.reduce(
                            (a, c) => a + c.quantity * c.price,
                            0
                        )}
                    </div>
                </li>
                <li>
                    <button
                        onClick={() => router.push("login?redirect=/shipping")}
                        className="primary-button w-full"
                    >
                        Check Out
                    </button>
                </li>
            </ul>
        </div>
    );
}
