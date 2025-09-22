#include <iostream>
#include <string>
#include <thread>
#include <asio.hpp>

int main()
{
    try
    {
        asio::io_context io;
        asio::ip::tcp::acceptor acceptor(io, {asio::ip::tcp::v4(), 9090});
        std::cout << "C++ stub listening on :9090" << std::endl;
        for (;;)
        {
            asio::ip::tcp::socket socket(io);
            acceptor.accept(socket);
            std::string resp =
                "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 2\r\n\r\nok";
            asio::write(socket, asio::buffer(resp));
        }
    }
    catch (std::exception &e)
    {
        std::cerr << "Error: " << e.what() << std::endl;
    }
    return 0;
}
