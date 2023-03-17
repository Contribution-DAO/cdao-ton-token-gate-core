import { Telegraf } from "telegraf"
import { v4 as uuidV4 } from 'uuid'
import { config as dotenvConfig } from 'dotenv'
import axios from "axios"

dotenvConfig()

const bot = new Telegraf(process.env.BOT_TOKEN as string)

bot.start((ctx) => {
  let message = `Hello world`
  ctx.reply(message)
})

bot.on('chat_join_request', async (ctx) => {
  let approved = false

  console.log("==============> new member joined")
  const chatId = ctx.chat?.id
  const userId = ctx.chatJoinRequest.user_chat_id

  try {
    if (chatId) {
      console.log(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId)
      const membership = await axios.get(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId, {
        headers: {
          Authorization: "Bearer " + process.env.API_SECRET,
        }
      })
  
      if (membership.data.isMinted && !membership.data.isJoined) {
        await ctx.approveChatJoinRequest(userId)
        approved = true
        await axios.post(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId + "/mark_joined", {
          joined: true,
        }, {
          headers: {
            Authorization: "Bearer " + process.env.API_SECRET,
          }
        })
      } else {
        await ctx.declineChatJoinRequest(userId)
        // return await ctx.telegram.sendMessage(chatId, "Rejected " + ctx.chatJoinRequest.from.first_name + "! Please tell him to join again after minting his SBT.")
      }
      // return await ctx.telegram.sendMessage(chatId, "Pending " + ctx.chatJoinRequest.from.first_name + " (ID: " + ctx.chatJoinRequest.user_chat_id + ")")
    }
  } catch (err) {
    if (!approved) {
      await ctx.declineChatJoinRequest(userId)
    }
  }
})

// bot.on('chat_member', async (ctx) => {
//   let approved = false;

//   console.log("==============> new member updated")
//   const chatId = ctx.chat?.id
//   const userId = ctx.chatMember.new_chat_member.user.id

//   try {
//     if (chatId) {
//       console.log(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId)
//       const membership = await axios.get(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId, {
//         headers: {
//           Authorization: "Bearer " + process.env.API_SECRET,
//         }
//       })
  
//       if (membership.data.isMinted && !membership.data.isJoined) {
//         await ctx.approveChatJoinRequest(userId)
//         approved = true
//         await axios.post(process.env.API_HOST + "/internal/" + userId + "/telegram/groups/" + chatId + "/mark_joined", {
//           joined: true,
//         }, {
//           headers: {
//             Authorization: "Bearer " + process.env.API_SECRET,
//           }
//         })
//       } else {
//         await ctx.declineChatJoinRequest(userId)
//         // return await ctx.telegram.sendMessage(chatId, "Rejected " + ctx.chatJoinRequest.from.first_name + "! Please tell him to join again after minting his SBT.")
//       }
//       // return await ctx.telegram.sendMessage(chatId, "Pending " + ctx.chatJoinRequest.from.first_name + " (ID: " + ctx.chatJoinRequest.user_chat_id + ")")
//     }
//   } catch (err) {
//     if (!approved) {
//       await ctx.declineChatJoinRequest(userId)
//     }
//   }
// })

bot.on('left_chat_member', async (ctx) => {
  try {
    console.log("==============> member left")
    const member = ctx.message.left_chat_member
    const chatId = ctx.chat?.id
  
    await axios.post(process.env.API_HOST + "/internal/" + member.id + "/telegram/groups/" + chatId + "/mark_joined", {
      joined: false,
    }, {
      headers: {
        Authorization: "Bearer " + process.env.API_SECRET,
      }
    })
  } catch (err) {

  }

})

bot.on('new_chat_members', async (ctx) => {
  try {
    console.log("==============> new member added")
    const newMembers = ctx.message.new_chat_members
    const chatId = ctx.chat?.id
    const botId = (await bot.telegram.getMe()).id
  
    if (chatId) {
      for (const member of newMembers) {
        if (member.id === botId) {
          // Bot is a new member of the chat
          await bot.telegram.sendMessage(chatId, `Please input this group ID in Ton connect UI: ${chatId}`)
        } else {
          // New members have joined the chat
          const membership = await axios.get(process.env.API_HOST + "/internal/" + member.id + "/telegram/groups/" + chatId, {
            headers: {
              Authorization: "Bearer " + process.env.API_SECRET,
            }
          })

          if (membership.data.isMinted) {
            await bot.telegram.sendMessage(chatId, "Hello " + member.first_name /*+ " (ID: " + member.id + ")"*/)
          } else {
            await ctx.kickChatMember(member.id)
          }
        }
      }
  
      // return await ctx.telegram.sendMessage(chatId, "Hello " + newMembers[0].first_name + " (ID: " + newMembers[0].id + ")")
    }
  } catch (err) {

  }

});


bot.launch()

console.log("Bot launched")