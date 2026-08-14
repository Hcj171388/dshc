package com.deepseek.dsh_android

import android.os.Bundle
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine
import io.flutter.plugin.common.MethodChannel
import mobile.Harness

class MainActivity: FlutterActivity() {
    private var harness: Harness? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        val dataDir = filesDir.absolutePath
        harness = try {
            Harness(dataDir)
        } catch (e: Exception) {
            e.printStackTrace()
            null
        }
    }

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)

        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, "com.deepseek.harness/agent").apply {
            setMethodCallHandler { call, result ->
                when (call.method) {
                    "createSession" -> result.success(harness?.createSession() ?: "")
                    "listSessions" -> result.success(harness?.listSessions() ?: "[]")
                    "getSession" -> {
                        val id = call.argument<String>("id") ?: ""
                        result.success(harness?.getSession(id) ?: "null")
                    }
                    "deleteSession" -> {
                        val id = call.argument<String>("id") ?: ""
                        harness?.deleteSession(id)
                        result.success(null)
                    }
                    "archiveSession" -> {
                        val id = call.argument<String>("id") ?: ""
                        harness?.archiveSession(id)
                        result.success(null)
                    }
                    "updateSessionTitle" -> {
                        val id = call.argument<String>("id") ?: ""
                        val title = call.argument<String>("title") ?: ""
                        harness?.updateSessionTitle(id, title)
                        result.success(null)
                    }
                    "getEvents" -> {
                        val sessionId = call.argument<String>("sessionId") ?: ""
                        val afterId = call.argument<Long>("afterId") ?: 0L
                        val limit = call.argument<Long>("limit") ?: 100L
                        result.success(harness?.getEvents(sessionId, afterId, limit) ?: "[]")
                    }
                    "runAgent" -> {
                        val sessionId = call.argument<String>("sessionId") ?: ""
                        val prompt = call.argument<String>("prompt") ?: ""
                        result.success(harness?.runAgent(sessionId, prompt) ?: "")
                    }
                    "abortAgent" -> {
                        harness?.abortAgent()
                        result.success(null)
                    }
                    else -> result.notImplemented()
                }
            }
        }

        MethodChannel(flutterEngine.dartExecutor.binaryMessenger, "com.deepseek.harness/config").apply {
            setMethodCallHandler { call, result ->
                when (call.method) {
                    "getConfig" -> result.success(harness?.getConfig() ?: "{}")
                    "saveConfig" -> {
                        val configJson = call.argument<String>("config") ?: "{}"
                        harness?.saveConfig(configJson)
                        result.success(null)
                    }
                    else -> result.notImplemented()
                }
            }
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        harness?.close()
    }
}
