package com.xinou.volunteer.utils;

/**
 * Created by zhangbo on 15/12/25.
 * 获取保证唯一的递增的序列号字符串
 */
public class IdUtil {


    public static String getTimeMillisSequence() {
        long nanoTime = System.nanoTime();
        if (nanoTime < 0) {
            nanoTime = nanoTime + Long.MAX_VALUE + 1;
        }
        String nanoTimeStr = String.valueOf(nanoTime);
        String preFix = "";
        int difBit = String.valueOf(Long.MAX_VALUE).length() - nanoTimeStr.length();
        for (int i = 0; i < difBit; i++) {
            preFix = preFix + "0";
        }
        nanoTimeStr = preFix + nanoTimeStr;
        String timeMillisSequence = String.valueOf(System.currentTimeMillis()) + "-" + nanoTimeStr;
        return timeMillisSequence;
    }
}
