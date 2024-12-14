/*
 * Variant 5.1.1
 * version 1.5
 *
 * Multithreading crawler
 */

#include <iostream>
#include <filesystem>
#include <unordered_set>
#include <fstream>
#include <ctime>
#include <thread>
#include <vector>
#include <mutex>

using namespace std;
namespace fs = filesystem;

fs::path sourcePath = fs::path(R"(C:\Users\aleks\Downloads\todo)");  // fs::current_path() <- maybe
fs::path pathTo = sourcePath / "to";
fs::path pathFrom = sourcePath / "from";
unordered_set<string> all_files;  // All files by pathFrom
unordered_set<string> checked_files;  // All files after 'check_file()'
unordered_set<string> rest_files;  // Needs to check
int number_actives = 0;
vector<thread *> threads;
int n;
mutex m1, m2;


bool starts_with(const string &s, const string &a) {  // s == a + s[n:]?
    if (a.size() > s.size()) return false;
    for (int i = 0; i < a.size(); ++i) {
        if (s[i] != a[i]) return false;
    }
    return true;
}

bool ends_with(const string &s, const string &a) {  // s == s[:-n] + a?
    if (a.size() > s.size()) return false;
    for (int i = 0; i < a.size(); ++i) {
        if (s.rbegin()[i] != a.rbegin()[i]) return false;
    }
    return true;
}

int custom_find(const string &s, const string &a) {  // If a not in s -> -1. Else (s[i:n] == a) -> i
    for (int i = 0; i < s.size(); ++i) {
        for (int j = 0; j < a.size(); ++j) {
            if (s[i + j] != a[j]) break;
            if (j + 1 == a.size()) return i;
        }
    }
    return -1;
}

// Выбирает файл из rest_files, если ранее с ним не работали, начинает работать: находит все ссылки в нем.
// Далее копирует соответствующие документы в директорию(якобы скачивает) и добавляет названия в rest_files
void check_file() {
    while (true) {
        m1.lock();
        if (!rest_files.empty()) {
            ++number_actives;
            string source_filename = *rest_files.begin();
            rest_files.erase(source_filename);
            if (!checked_files.contains(source_filename)) {
                checked_files.emplace(source_filename);
                m1.unlock();
                ifstream inlet((pathTo / source_filename).string());
                string s;
                while (inlet >> s) {
                    if (ends_with(s, "<a")) {
                        inlet >> s;
                        if (starts_with(s, "href=\"")) {
                            int ind_end = custom_find(s, "\">");
                            if (ind_end != -1) {
                                string new_filename = s.substr(6, ind_end - 6);
                                if (starts_with(new_filename, "file://")) {
                                    new_filename = new_filename.substr(7, new_filename.size() - 7);
                                    lock_guard<mutex> lock(m2);
                                    if (!checked_files.contains(new_filename) && all_files.contains(new_filename) &&
                                        !rest_files.contains(new_filename)) {
                                        fs::copy(pathFrom / new_filename, pathTo / new_filename);
                                        rest_files.emplace(new_filename);
                                    }
                                }
                            }
                        }
                    }
                }
                inlet.close();
            } else m1.unlock();
            --number_actives;
        } else {
            m1.unlock();
            if (number_actives == 0) break;
        }
    }
}

int main() {
    string s;
    cin >> s >> n;

    unsigned int start_time = clock();
    for (const fs::directory_entry &dir: fs::directory_iterator(pathFrom)) {
        all_files.emplace(dir.path().filename().string());
    }

    if (fs::exists(pathTo)) fs::remove_all(pathTo);
    fs::create_directory(pathTo);

    fs::copy(pathFrom / s, pathTo / s);
    rest_files.emplace(s);

    for (int i = 0; i < n; ++i) threads.emplace_back(new thread(check_file));

    for (thread *x: threads) x->join();

    fs::remove_all(pathTo);
    cout << checked_files.size() << ' ' << clock() - start_time << endl;
    return 0;
}
/* s = "0.html"
 * number_threads | operation_time
 *      1         |      3406
 *      2         |      2148
 *      3         |      1723
 *      4         |      1639
 *      5         |      1572
 *      6         |      1475
 *      7         |      1553
 */
