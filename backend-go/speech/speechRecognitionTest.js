const fs = require('fs');
const AipSpeechClient = require("baidu-aip-sdk").speech;
const HttpClient = require("baidu-aip-sdk").HttpClient;

// 设置APPID/AK/SK（使用提供的测试账号）
const APP_ID = "120154300";
const API_KEY = "ANjuQLa2GjY8pWvJXbjf7LzI";
const SECRET_KEY = "107wk6X4Tak7gmb5nYnEBgXMdgGbHCsf";

// 新建一个语音识别客户端实例
const client = new AipSpeechClient(APP_ID, API_KEY, SECRET_KEY);

// 配置请求参数，设置超时时间为5秒
HttpClient.setRequestOptions({ timeout: 5000 });

/**
 * 语音识别函数
 * @param {string} filePath - 语音文件路径
 * @param {string} format - 语音文件格式，如pcm、wav、amr
 * @param {number} rate - 采样率，16000或8000
 * @param {object} options - 可选参数
 */
async function recognizeSpeech(filePath, format = 'pcm', rate = 16000, options = {}) {
    try {
        // 读取语音文件并转换为Buffer
        const voice = fs.readFileSync(filePath);
        const voiceBuffer = new Buffer(voice);
        
        // 调用百度语音识别API
        const result = await client.recognize(voiceBuffer, format, rate, options);
        
        console.log('识别结果:', JSON.stringify(result, null, 2));
        
        // 处理识别结果
        if (result.err_no === 0) {
            console.log('识别成功，内容:', result.result.join(','));
            return result.result;
        } else {
            console.error('识别失败:', result.err_msg);
            return null;
        }
    } catch (err) {
        console.error('识别过程出错:', err);
        return null;
    }
}

// 测试语音识别
// 请将下面的文件路径替换为你的语音文件实际路径
const testFilePath = './16k.pcm'; // 示例：PCM格式的语音文件

// 调用识别函数，使用普通话模型（1537）
recognizeSpeech(testFilePath, 'pcm', 16000, {
    dev_pid: 1537, // 普通话(纯中文识别)
    cuid: 'test-device-' + Math.random().toString(36).substr(2, 9) // 生成随机设备ID
});
