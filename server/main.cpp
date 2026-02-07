#include <iostream>
#include <cstring>
#include <arpa/inet.h>
#include <unistd.h>
#include <stdexcept>

void handleClient(const int serverSocket);
std::vector<std::string> parsePacket(std::string packet);
std::string login() {
    return "Enter your username and password: \n";
}
std::unordered_map<std::string, float> accounts = { {"hello", 100.0f} };

int main() {
    int serverSocket = socket(AF_INET, SOCK_DGRAM, 0);

    sockaddr_in serverAddr;
    serverAddr.sin_family = AF_INET;
    serverAddr.sin_addr.s_addr = INADDR_ANY; // Listen on any IP
    serverAddr.sin_port = htons(2222);       // Port 2222

    // Bind the socket to the port
    bind(serverSocket, reinterpret_cast<sockaddr *>(&serverAddr), sizeof(serverAddr));

    std::cout << "C++ Server running on port 2222..." << std::endl;

    while (true) {
        handleClient(serverSocket);
    }

    close(serverSocket);
    return 0;
}

void handleClient(const int serverSocket) {
    char buffer[1000];
    sockaddr_in clientAddr;
    socklen_t clientLen = sizeof(clientAddr);

    // 1. Receive
    int bytesReceived = recvfrom(serverSocket, buffer, 1000, 0,
                                 reinterpret_cast<struct sockaddr *>(&clientAddr), &clientLen);
    
    if (bytesReceived < 0) return; // Network error, just return
    
    buffer[bytesReceived] = '\0'; 
    std::string request(buffer);
    std::cout << "Received: " << request << std::endl;

    std::string reply = "ERROR"; // Default reply

    try {
        std::vector<std::string> data = parsePacket(request);

        if (data.empty()) {
            throw std::runtime_error("Empty packet parsed");
        }

        // 2. Logic
        if (data[0] == "LOGIN") {
            if (data.size() < 3) {
                throw std::runtime_error("LOGIN requires 3 arguments");
            }

            std::string username = data[1];
            std::string passStr = data[2];

            if (accounts.find(username) != accounts.end()) {
                
                float receivedVal = 0.0f;
                size_t processedChars = 0;
                
                // std::stof can throw invalid_argument or out_of_range
                receivedVal = std::stof(passStr, &processedChars);

                if (std::abs(accounts[username] - receivedVal) <= 0.001f) {
                    reply = "SUCCESS";
                } else {
                    reply = "FAIL: Incorrect value";
                }
            } else {
                reply = "FAIL: User not found";
            }
        } else {
             reply = "UNKNOWN_COMMAND";
        }

    } catch (const std::invalid_argument& e) {
        std::cerr << " [!] Client sent text where number was expected." << std::endl;
        reply = "ERROR: Number format invalid";
    } catch (const std::out_of_range& e) {
        std::cerr << " [!] Client sent a number too big for float/int." << std::endl;
        reply = "ERROR: Number too large";
    } catch (const std::exception& e) {
        // Catch-all for any other logic errors (like vector index out of bounds)
        std::cerr << " [!] Standard Error: " << e.what() << std::endl;
        reply = "ERROR: Bad Request";
    } catch (...) {
        std::cerr << " [!] Unknown crash prevented." << std::endl;
    }

    // 3. Reply (Always reply, even if it's an error message)
    sendto(serverSocket, reply.c_str(), reply.length(), 0,
           reinterpret_cast<sockaddr *>(&clientAddr), clientLen);
}

std::vector<std::string> parsePacket(std::string packet) {
    std::vector<std::string> parts;
    int cursor = 0; // Points to where we are currently reading

    while (cursor < packet.length()) {
        // 1. Find the next colon (separates length from data)
        size_t delimiterPos = packet.find(':', cursor);
        
        if (delimiterPos == std::string::npos) break; // formatting error or end

        // 2. Read the number string (e.g., "5") and convert to int
        std::string lengthStr = packet.substr(cursor, delimiterPos - cursor);
        int length = std::stoi(lengthStr);

        // 3. Jump over the colon to the start of the data
        int dataStart = delimiterPos + 1;

        // 4. Extract the exact number of characters
        std::string data = packet.substr(dataStart, length);
        parts.push_back(data);

        // 5. Move cursor to the start of the NEXT chunk
        cursor = dataStart + length;
    }
    return parts;
}