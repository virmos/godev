import Image from "next/image";

export default function FeatureSection() {
    return (
        <section id="feature">
            {/* <!-- Feature Container --> */}
            <div className="relative container flex flex-col max-w-6xl mx-auto my-32 px-6 text-gray-900 md:flex-row md:px-0 font-sans">
                {/* <!-- Image --> */}
                <Image
                    src="/images/desktop/image-interactive.jpg"
                    alt=""
                    width="600"
                    height="500"
                ></Image>

                {/* <!-- Text Container --> */}
                <div className="top-32 pr-0 bg-white md:absolute md:right-0 md:py-20 md:pl-20">
                    <h2 className="max-w-lg mt-10 mb-6 font-sans text-4xl text-center text-gray-900 uppercase md:text-5xl md:mt-0 md:text-left">
                        The leader in electronic commerce
                    </h2>

                    <p className="max-w-md text-center md:text-left">
                        Founded in 1994, Shopee has earned a reputation as a
                        disruptor of well-established industries through
                        technological innovation and aggressive reinvestment
                        of profits into capital expenditures. As of 2023, it
                        is the worlds largest online retailer and
                        marketplace.
                    </p>
                </div>
            </div>
        </section>
    );
}
